package domain_model

type Passage struct {
    Id int
    Reference PassageReference
    Text string
    Owner int
    ReviewState *PassageReview
}

func NewPassage(reference PassageReference, text string, owner int) Passage {
    return Passage {
        Reference: reference,
        Text: text,
        Owner: owner,
    }
}

func (p *Passage) SetReference(reference PassageReference) {
    p.Reference = reference
}

func (p *Passage) SetText(text string) {
    p.Text = text
}

func (p *Passage) SetReviewState(reviewState *PassageReview) {
    p.ReviewState = reviewState
}

func (p *Passage) Review(grade ReviewGrade, timestamp ReviewTimestamp) {
    if p.ReviewState == nil {
        nextReview := FirstReview(grade, timestamp)
        p.ReviewState = &nextReview
    } else {
        nextReview := p.ReviewState.NextAfterReview(grade, timestamp)
        p.ReviewState = &nextReview
    }
}
