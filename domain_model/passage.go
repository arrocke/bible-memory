package domain_model

type Passage struct {
    new bool

    id int
    reference PassageReference
    text string
    owner int
    reviewState *PassageReview
}

func NewPassage(reference PassageReference, text string, owner int) Passage {
    return Passage {
        reference: reference,
        text: text,
        owner: owner,
    }
}

func (p *Passage) SetReference(reference PassageReference) {
    p.reference = reference
}

func (p *Passage) SetText(text string) {
    p.text = text
}

func (p *Passage) SetReviewState(interval ReviewInterval, nextReview ReviewTimestamp) {
    nextReviewState := p.reviewState.Overwrite(interval, nextReview)
    p.reviewState = &nextReviewState
}

func (p *Passage) Review(grade ReviewGrade, timestamp ReviewTimestamp) {
    if p.reviewState == nil {
        nextReview := FirstReview(grade, timestamp)
        p.reviewState = &nextReview
    } else {
        nextReview := p.reviewState.NextAfterReview(grade, timestamp)
        p.reviewState = &nextReview
    }
}
