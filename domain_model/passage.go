package domain_model

import (
	"strconv"
)

type PassageId int

func ParsePassageId(idstr string) (PassageId, error) {
    id, err := strconv.ParseUint(idstr, 10, 32)
    if err != nil {
        return 0, err
    }
    return PassageId(id), nil
}

type Passage struct {
    new bool

    Id Id
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

func (p *Passage) SetReviewState(interval ReviewInterval, nextReview ReviewTimestamp) {
    nextReviewState := p.ReviewState.Overwrite(interval, nextReview)
    p.ReviewState = &nextReviewState
}

func (p *Passage) Review(grade ReviewGrade, timestamp ReviewTimestamp) {
    nextReview := p.ReviewState.NextAfterReview(grade, timestamp)
    p.ReviewState = &nextReview
}
