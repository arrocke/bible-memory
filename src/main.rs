mod passage;
mod routes;

use axum::{serve, Router};
use sqlx::postgres::PgPoolOptions;
use std::net::Ipv4Addr;
use tokio::net::TcpListener;

use routes::AppState;

#[tokio::main]
async fn main() {
    let db_pool = PgPoolOptions::new()
        .max_connections(5)
        .connect("postgres://adrian:adrian@localhost:5432/bible-memory")
        .await
        .unwrap();

    let app = Router::new()
        .merge(routes::get_passages::route())
        .merge(routes::post_passage::route())
        .merge(routes::get_new_passage::route())
        .merge(routes::get_passage_review::route())
        .with_state(AppState { db_pool });
    println!("See example: http://127.0.0.1:8080/example");

    let listener = TcpListener::bind((Ipv4Addr::LOCALHOST, 8080))
        .await
        .unwrap();
    serve(listener, app.into_make_service()).await.unwrap();
}
