use askama::Template;
use askama_axum::IntoResponse;
use axum::{
    extract::{Path, State},
    http::StatusCode,
    routing::get,
    Router,
};

use crate::passage::{Passage, PassageReference};
use crate::routes::{AppState, DbPool, ErrorTemplate, NotFoundTemplate};

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
    let Ok(result) = query_passage(&db_pool, passage_id).await else {
        return (
            StatusCode::INTERNAL_SERVER_ERROR,
            ErrorTemplate {
                error_message: String::from("Unknown error occurred."),
            },
        )
            .into_response();
    };
    let Some(passage) = result else {
        return (
            StatusCode::NOT_FOUND,
            NotFoundTemplate {
                error_message: String::from("Couldn't find passage to review."),
            },
        )
            .into_response();
    };
    (ReviewTemplate { passage }).into_response()
}

pub fn route() -> Router<AppState> {
    Router::new().route("/passages/:passage_id/review", get(handler))
}
