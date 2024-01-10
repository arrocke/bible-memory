use crate::passage::{Passage, PassageReference};
use askama::Template;
use axum::{extract::State, response::IntoResponse, routing::get, Router};

use crate::routes::{AppState, DbPool, ErrorResponse};

#[derive(Template)]
#[template(path = "index.html")]
struct IndexTemplate {
    passages: Vec<Passage>,
}

async fn query_passages(db_pool: &DbPool) -> Result<Vec<Passage>, sqlx::Error> {
    sqlx::query!(r#"SELECT * FROM passage"#)
        .fetch_all(db_pool)
        .await
        .map(|rows| {
            rows.iter()
                .map(|row| Passage {
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
                .collect()
        })
}

async fn handler(State(db_pool): State<DbPool>) -> Result<impl IntoResponse, ErrorResponse> {
    let template = IndexTemplate {
        passages: query_passages(&db_pool).await?,
    };
    Ok(template)
}

pub fn route() -> Router<AppState> {
    Router::<AppState>::new().route("/", get(handler))
}
