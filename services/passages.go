package services

import (
	"main/db"
	"main/domain_model"
	"time"
)

type PassagesService struct {
	passageRepo db.PassageRepo
}

func CreatePassagesService(passageRepo db.PassageRepo) PassagesService {
	return PassagesService{passageRepo}
}

type CreatePassageRequest struct {
	Reference string
	Text      string
	UserId    int
}

func (service *PassagesService) Create(request CreatePassageRequest) (int, error) {
	reference, err := domain_model.ParsePassageReference(request.Reference)
	if err != nil {
		return 0, err
	}
	passage := domain_model.NewPassage(reference, request.Text, request.UserId)

	if err := service.passageRepo.Create(&passage); err != nil {
		return 0, nil
	}

	return passage.Id, nil
}

type UpdatePassageRequest struct {
	Id        int
	Reference string
	Text      string
	Interval  *int
	ReviewAt  *time.Time
}

func (service *PassagesService) Update(request UpdatePassageRequest) error {
	passage, err := service.passageRepo.Get(request.Id)
	if err != nil {
		return err
	}

	reference, err := domain_model.ParsePassageReference(request.Reference)
	if err != nil {
		return err
	}

	var nextReview *domain_model.PassageReview
	if request.Interval != nil && request.ReviewAt != nil {
		interval, err := domain_model.NewReviewInterval(*request.Interval)
		if err != nil {
			return err
		}

		reviewAt := domain_model.NewReviewTimestamp(*request.ReviewAt)

		reviewState := passage.ReviewState.Overwrite(interval, reviewAt)
		nextReview = &reviewState
	}

	passage.SetReference(reference)
	passage.SetText(request.Text)
	passage.SetReviewState(nextReview)

	return service.passageRepo.Update(&passage)
}

type ReviewPassageRequest struct {
	Id    int
	Grade int
    Tz int
}

func (service *PassagesService) Review(request ReviewPassageRequest) error {
	passage, err := service.passageRepo.Get(request.Id)
	if err != nil {
		return err
	}

    grade, err := domain_model.NewReviewGrade(request.Grade)
    if err != nil {
        return err
    }

    timestamp := domain_model.NewReviewTimestampForToday(request.Tz)

    passage.Review(grade, timestamp) 
    
    err = service.passageRepo.Update(&passage)
    if err != nil {
        return err
    }

    return nil
}
