use std::net::Ipv4Addr;

use askama::Template;
use axum::{
    extract::{Form, FromRef, Path, State},
    response::{IntoResponse, Redirect},
    routing::{get, post},
    serve, Router,
};
use serde::Deserialize;
use tokio::net::TcpListener;

#[derive(Debug, Clone)]
pub struct Passage {
    id: u32,
    reference: String,
    level: u32,
}

#[derive(Template)]
#[template(path = "index.html")]
struct IndexTemplate {
    passages: Vec<Passage>,
}

async fn get_index(State(passages): State<Vec<Passage>>) -> impl IntoResponse {
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
    State(passages): State<Vec<Passage>>,
    Form(form): Form<NewPassageForm>,
) -> impl IntoResponse {
    let id = match passages.iter().map(|p| p.id).max() {
        Some(id) => id + 1,
        None => 1,
    };
    println!("New passage {} {}", id, form.reference);
    // let mut mut_passages = &passages;
    // mut_passages.push(Passage {
    //     id,
    //     reference: form.reference,
    //     level: 0,
    // });
    Redirect::to("/")
}

#[derive(Template)]
#[template(path = "review.html")]
struct ReviewTemplate {
    passage: Passage,
}

async fn get_review(
    State(passages): State<Vec<Passage>>,
    Path(resource_id): Path<u32>,
) -> impl IntoResponse {
    let Some(passage) = passages.iter().find(|&p| p.id == resource_id) else {
        panic!("Not found")
    };
    ReviewTemplate {
        passage: passage.clone(),
    }
}

#[derive(Clone, FromRef)]
struct AppState {
    passages: Vec<Passage>,
}

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/", get(get_index))
        .route("/passages", post(post_new_passage))
        .route("/passages/new", get(get_new_passage))
        .route("/passages/:passage_id/review", get(get_review))
        // Create the application state
        .with_state(AppState {
            passages: Vec::from([
                Passage {
                    id: 1,
                    reference: String::from("Genesis 1:1-5"),
                    level: 1,
                },
                Passage {
                    id: 2,
                    reference: String::from("Genesis 2:5-10"),
                    level: 2,
                },
            ]),
        });
    println!("See example: http://127.0.0.1:8080/example");

    let listener = TcpListener::bind((Ipv4Addr::LOCALHOST, 8080))
        .await
        .unwrap();
    serve(listener, app.into_make_service()).await.unwrap();
}
