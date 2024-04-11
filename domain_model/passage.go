package domain_model

import "math/rand"

type PassageProps struct {
    Id int
    Reference PassageReference
    Text string
    Owner int
    ReviewState *PassageReview
}

type Passage struct {
    props PassageProps
    isNew bool
}

func PassageFactory(props PassageProps) Passage {
    return Passage {
        props: props,
        isNew: false,
    }
}

type NewPassageProps struct {
    Reference PassageReference
    Text string
    Owner int
}

func NewPassage(props NewPassageProps) Passage {
    id := rand.Int()
    return Passage{
        props: PassageProps {
            Id: id,
            Reference: props.Reference,
            Text: props.Text,
            Owner: props.Owner,
        },
        isNew: true,
    }
}

func (u *Passage) IsNew() bool {
    return u.isNew
}

func (u *Passage) Id() int {
    return u.props.Id
}

func (u *Passage) Props() PassageProps {
    return u.props
}

func (p *Passage) SetReference(reference PassageReference) {
    p.props.Reference = reference
}

func (p *Passage) SetText(text string) {
    p.props.Text = text
}

func (p *Passage) SetReviewState(reviewState *PassageReview) {
    p.props.ReviewState = reviewState
}

func (p *Passage) Review(grade ReviewGrade, timestamp ReviewTimestamp) {
    if p.props.ReviewState == nil {
        nextReview := FirstReview(grade, timestamp)
        p.props.ReviewState = &nextReview
    } else {
        nextReview := p.props.ReviewState.NextAfterReview(grade, timestamp)
        p.props.ReviewState = &nextReview
    }
}
