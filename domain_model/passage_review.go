package domain_model

type PassageReview struct {
	Interval   ReviewInterval
	NextReview ReviewTimestamp
	ReviewedAt *ReviewTimestamp
}

func FirstReview(grade ReviewGrade, timestamp ReviewTimestamp) PassageReview {
    interval := FirstReviewInterval(grade)
    return PassageReview {
        Interval: interval,
        NextReview: timestamp.AddInterval(interval),
        ReviewedAt: &timestamp,
    }
}

func (current PassageReview) NextAfterReview(grade ReviewGrade, timestamp ReviewTimestamp) PassageReview {
    actualInterval := current.Interval
    if current.ReviewedAt != nil {
	    actualInterval =  current.ReviewedAt.IntervalToDate(timestamp)
    }

    nextInterval := current.Interval.NextAfterReview(grade, actualInterval)

    return PassageReview {
        Interval: nextInterval,
        NextReview: timestamp.AddInterval(nextInterval),
        ReviewedAt: &timestamp,
    }
}

func (r *PassageReview) Overwrite(interval ReviewInterval, nextReview ReviewTimestamp) PassageReview {
    var reviewedAt *ReviewTimestamp
    if r != nil {
        reviewedAt = r.ReviewedAt
    }

	return PassageReview{
		Interval:   interval,
		ReviewedAt: reviewedAt,
		NextReview: nextReview,
	}
}

