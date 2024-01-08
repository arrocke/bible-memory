use askama_axum::IntoResponse;
use axum::{extract::State, response::Redirect, routing::post, Form, Router};
use serde::Deserialize;

use crate::passage::PassageReference;

use crate::routes::{AppState, DbPool};

#[derive(Deserialize)]
struct NewPassageForm {
    reference: String,
}

struct Passage {
    reference: PassageReference,
}

async fn insert_passage(db_pool: &DbPool, passage: &Passage) -> Result<(), sqlx::Error> {
    sqlx::query!(
        r#"INSERT INTO passage (book, start_chapter, start_verse, end_chapter, end_verse) VALUES ($1, $2, $3, $4, $5)"#,
        passage.reference.book,
        passage.reference.start_chapter,
        passage.reference.start_verse,
        passage.reference.end_chapter,
        passage.reference.end_verse
    ).execute(db_pool).await.and(Ok(()))
}

async fn handler(
    State(db_pool): State<DbPool>,
    Form(form): Form<NewPassageForm>,
) -> impl IntoResponse {
    let passage = Passage {
        reference: form.reference.parse::<PassageReference>().unwrap(),
    };
    insert_passage(&db_pool, &passage).await.unwrap();
    Redirect::to("/")
}

pub fn route() -> Router<AppState> {
    Router::new().route("/passages", post(handler))
}
