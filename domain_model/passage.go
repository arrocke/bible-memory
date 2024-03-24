package domain_model

import "time"

type Passage struct {
    new bool

    Id int
    Reference PassageReference
    Text string
    Owner int
    ReviewState *PassageReview
}

func (p Passage) IsNew() bool {
    return p.new
}

func NewPassage(reference PassageReference, text string, owner int) Passage {
    return Passage {
        new: true,

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

func (p *Passage) OverrideReviewState(interval ReviewInterval, nextReview ReviewTimestamp) {
    nextReviewState := p.ReviewState.Update(interval, nextReview)
    p.ReviewState = &nextReviewState
}

func (p *Passage) Review(grade ReviewGrade, timestamp time.Time) {
    nextReview := p.ReviewState.Review(grade, ReviewTimestamp(timestamp))
    p.ReviewState = &nextReview
}
