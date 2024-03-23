package domain_model

import "time"

type Passage struct {
    new bool

    Id uint
    Reference PassageReference
    Text string
    Owner uint
    Review *PassageReview
}

func (p Passage) IsNew() bool {
    return p.new
}

func NewPassage(reference PassageReference, text string, owner uint) Passage {
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

func (p *Passage) SetReview(review PassageReview) {
    p.Review = &review
}

func (p *Passage) DoReview(grade uint, timestamp time.Time) error {
    nextReview, err := p.Review.Review(grade, timestamp)
    if err != nil {
        return err
    }
    p.Review = &nextReview
    return nil
}
