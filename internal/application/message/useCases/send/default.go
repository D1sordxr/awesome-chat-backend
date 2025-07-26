package send

import (
	"awesome-chat/internal/application/message/dto"
	"awesome-chat/internal/domain/core/message/ports"
	outboxPorts "awesome-chat/internal/domain/core/shared/outbox/ports"
	"awesome-chat/internal/domain/core/shared/outbox/vo"
	sharedPorts "awesome-chat/internal/domain/core/shared/ports"
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
)

const (
	workersCount = 100
	jobsBuff     = 100
)

type work struct {
	ctx     context.Context
	req     dto.SendRequest
	resChan chan error
}

type UseCase struct {
	entityCreator ports.EntityCreator
	outboxCreator outboxPorts.EntityCreator
	txManager     sharedPorts.TransactionManager
	msgRepo       ports.Repository
	outboxRepo    outboxPorts.Repository
	wg            sync.WaitGroup
	jobs          chan work
}

func (uc *UseCase) worker(id int) { // TODO: log with worker id
	defer uc.wg.Done()

	for job := range uc.jobs {
		select {
		case <-job.ctx.Done():
			job.resChan <- job.ctx.Err()
			continue
		default:
			func() {
				ctx, err := uc.txManager.BeginAndInjectTx(job.ctx)
				if err != nil {
					job.resChan <- err
					return
				}
				defer func() { _ = uc.txManager.RollbackTx(ctx) }()

				entity := uc.entityCreator.Do(job.req.UserID, job.req.ChatID, job.req.Content)
				if err = uc.msgRepo.Save(ctx, entity); err != nil {
					job.resChan <- err
					return
				}

				msgPayload, err := json.Marshal(entity)
				if err != nil {
					job.resChan <- err
					return
				}

				outbox := uc.outboxCreator.Do(
					uuid.New(),
					vo.MessageEntity,
					msgPayload,
				)
				if err = uc.outboxRepo.Save(ctx, outbox); err != nil {
					job.resChan <- err
					return
				}

				if err = uc.txManager.CommitTx(ctx); err != nil {
					job.resChan <- err
					return
				}

				job.resChan <- nil
			}()
		}
	}
}

func NewUseCase(
	entityCreator ports.EntityCreator,
	outboxCreator outboxPorts.EntityCreator,
	txManager sharedPorts.TransactionManager,
	msgRepo ports.Repository,
	outboxRepo outboxPorts.Repository,
) *UseCase {
	uc := &UseCase{
		entityCreator: entityCreator,
		outboxCreator: outboxCreator,
		txManager:     txManager,
		msgRepo:       msgRepo,
		outboxRepo:    outboxRepo,
		wg:            sync.WaitGroup{},
		jobs:          make(chan work, jobsBuff),
	}

	uc.wg.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		go uc.worker(i)
	}

	return uc
}

func (uc *UseCase) Execute(ctx context.Context, req dto.SendRequest) error {
	resChan := make(chan error, 1)
	w := work{
		ctx:     ctx,
		req:     req,
		resChan: resChan,
	}
	select {
	case uc.jobs <- w:
	case <-ctx.Done():
		return ctx.Err()
	}

	select {
	case err := <-resChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (uc *UseCase) Start(ctx context.Context) error { return nil }

func (uc *UseCase) Shutdown(ctx context.Context) error {
	close(uc.jobs)

	done := make(chan struct{})
	go func() {
		uc.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
