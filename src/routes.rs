pub mod get_new_passage;
pub mod get_passage_review;
pub mod get_passages;
pub mod post_passage;

use askama::Template;
use axum::{body::Body, extract::FromRef, http::Response, response::IntoResponse};
use sqlx::{postgres::Postgres, Pool};

pub type DbPool = Pool<Postgres>;

#[derive(Clone, FromRef)]
pub struct AppState {
    pub db_pool: DbPool,
}

pub enum ErrorResponse {
    NotFound { resource: String },
    ServerError,
}

impl From<sqlx::Error> for ErrorResponse {
    fn from(_value: sqlx::Error) -> Self {
        Self::ServerError
    }
}

#[derive(Template)]
#[template(path = "not-found.html")]
struct NotFoundTemplate {
    error_message: String,
}
#[derive(Template)]
#[template(path = "error.html")]
struct ErrorTemplate;

impl IntoResponse for ErrorResponse {
    fn into_response(self) -> Response<Body> {
        match self {
            ErrorResponse::NotFound { resource } => (NotFoundTemplate {
                error_message: format!("Could not find {}", resource),
            })
            .into_response(),
            ErrorResponse::ServerError => ErrorTemplate.into_response(),
        }
    }
}
