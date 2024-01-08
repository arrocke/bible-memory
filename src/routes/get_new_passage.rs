use askama::Template;
use axum::{response::IntoResponse, routing::get, Router};

use crate::routes::AppState;

#[derive(Template)]
#[template(path = "new-passage.html")]
struct NewPassageTemplate {}

async fn handler() -> impl IntoResponse {
    NewPassageTemplate {}
}

pub fn route() -> Router<AppState> {
    Router::new().route("/passages/new", get(handler))
}
