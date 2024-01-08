use askama::Template;
use askama_axum::IntoResponse;
use axum::{
    extract::{Path, State},
    routing::get,
    Router,
};

use crate::passage::{Passage, PassageReference};
use crate::routes::{AppState, DbPool};

#[derive(Template)]
#[template(path = "review.html")]
struct ReviewTemplate {
    passage: Passage,
}

async fn query_passage(db_pool: &DbPool, passage_id: i32) -> Result<Option<Passage>, sqlx::Error> {
    sqlx::query!(r#"SELECT * FROM passage WHERE id = $1"#, passage_id)
        .fetch_optional(db_pool)
        .await
        .map(|row| {
            row.and_then(|row| {
                Some(Passage {
                    id: row.id,
                    reference: PassageReference {
                        book: row.book.clone(),
                        start_chapter: row.start_chapter,
                        start_verse: row.start_verse,
                        end_chapter: row.end_chapter,
                        end_verse: row.end_verse,
                    },
                    level: 0,
                })
            })
        })
}

async fn handler(State(db_pool): State<DbPool>, Path(passage_id): Path<i32>) -> impl IntoResponse {
    let passage = query_passage(&db_pool, passage_id).await.unwrap().unwrap();
    ReviewTemplate { passage }
}

pub fn route() -> Router<AppState> {
    Router::new().route("/passages/:passage_id/review", get(handler))
}
