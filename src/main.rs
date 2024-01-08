use std::net::Ipv4Addr;

use askama::Template;
use axum::{
    extract::{Form, FromRef, Path, State},
    response::{IntoResponse, Redirect},
    routing::{get, post},
    serve, Router,
};
use serde::Deserialize;
use sqlx::{
    postgres::{PgPoolOptions, Postgres},
    Pool,
};
use tokio::net::TcpListener;

mod passage;

use passage::{Passage, PassageReference};

#[derive(Template)]
#[template(path = "index.html")]
struct IndexTemplate {
    passages: Vec<Passage>,
}

async fn get_index(State(db_pool): State<Pool<Postgres>>) -> impl IntoResponse {
    let passages: Vec<Passage> = sqlx::query!(r#"SELECT * FROM passage"#)
        .fetch_all(&db_pool)
        .await
        .unwrap()
        .iter()
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
        .collect();

    IndexTemplate { passages }
}

#[derive(Template)]
#[template(path = "new-passage.html")]
struct NewPassageTemplate {}

async fn get_new_passage() -> impl IntoResponse {
    NewPassageTemplate {}
}

#[derive(Deserialize)]
struct NewPassageForm {
    reference: String,
}

async fn post_new_passage(
    State(db_pool): State<Pool<Postgres>>,
    Form(form): Form<NewPassageForm>,
) -> impl IntoResponse {
    let reference = form.reference.parse::<PassageReference>().unwrap();
    sqlx::query!(
        r#"INSERT INTO passage (book, start_chapter, start_verse, end_chapter, end_verse) VALUES ($1, $2, $3, $4, $5)"#,
        reference.book,
        reference.start_chapter,
        reference.start_verse,
        reference.end_chapter,
        reference.end_verse
    ).execute(&db_pool).await.unwrap();
    Redirect::to("/")
}

#[derive(Template)]
#[template(path = "review.html")]
struct ReviewTemplate {
    passage: Passage,
}

async fn get_review(
    State(db_pool): State<Pool<Postgres>>,
    Path(resource_id): Path<i32>,
) -> impl IntoResponse {
    let row = sqlx::query!(r#"SELECT * FROM passage WHERE id = $1"#, resource_id)
        .fetch_one(&db_pool)
        .await
        .unwrap();
    ReviewTemplate {
        passage: Passage {
            id: row.id,
            reference: PassageReference {
                book: row.book.clone(),
                start_chapter: row.start_chapter,
                start_verse: row.start_verse,
                end_chapter: row.end_chapter,
                end_verse: row.end_verse,
            },
            level: 0,
        },
    }
}

#[derive(Clone, FromRef)]
struct AppState {
    db_pool: Pool<Postgres>,
}

#[tokio::main]
async fn main() {
    let db_pool = PgPoolOptions::new()
        .max_connections(5)
        .connect("postgres://adrian:adrian@localhost:5432/bible-memory")
        .await
        .unwrap();

    let app = Router::new()
        .route("/", get(get_index))
        .route("/passages", post(post_new_passage))
        .route("/passages/new", get(get_new_passage))
        .route("/passages/:passage_id/review", get(get_review))
        .with_state(AppState { db_pool });
    println!("See example: http://127.0.0.1:8080/example");

    let listener = TcpListener::bind((Ipv4Addr::LOCALHOST, 8080))
        .await
        .unwrap();
    serve(listener, app.into_make_service()).await.unwrap();
}
