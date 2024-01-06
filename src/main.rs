use std::net::Ipv4Addr;

use axum::{
    extract::{FromRef, Path, State},
    response::IntoResponse,
    routing::get,
    serve, Router,
};
use axum_template::{engine::Engine, Key, RenderHtml};
use handlebars::Handlebars;
use serde::Serialize;
use tokio::net::TcpListener;

// Type alias for our engine. For this example, we are using Handlebars
type AppEngine = Engine<Handlebars<'static>>;

#[derive(Debug, Serialize, Clone)]
pub struct Passage {
    id: u32,
    reference: String,
    level: u32,
    // review_date: std::time::Date
}

#[derive(Debug, Serialize)]
struct IndexData {
    title: String,
    passages: Vec<Passage>,
}

async fn get_index(
    engine: AppEngine,
    State(passages): State<Vec<Passage>>,
    Key(key): Key,
) -> impl IntoResponse {
    RenderHtml(
        key,
        engine,
        IndexData {
            title: String::from("Bible Memory"),
            passages,
        },
    )
}

#[derive(Debug, Serialize)]
struct ReviewData {
    title: String,
    passage: Passage,
}

async fn get_review(
    engine: AppEngine,
    State(passages): State<Vec<Passage>>,
    Key(key): Key,
    Path(resource_id): Path<u32>,
) -> impl IntoResponse {
    let Some(passage) = passages.iter().find(|&p| p.id == resource_id) else {
        panic!("Not found")
    };
    RenderHtml(
        key,
        engine,
        ReviewData {
            title: format!("Review {} | Bible Memory", passage.reference),
            passage: passage.clone(),
        },
    )
}

#[derive(Clone, FromRef)]
struct AppState {
    engine: AppEngine,
    passages: Vec<Passage>,
}

#[tokio::main]
async fn main() {
    let mut hbs = Handlebars::new();
    hbs.register_template_file("head", "./src/templates/partials/head.hbs")
        .unwrap();
    hbs.register_template_file("body", "./src/templates/partials/body.hbs")
        .unwrap();
    hbs.register_template_file("/", "./src/templates/pages/index.hbs")
        .unwrap();
    hbs.register_template_file("/:passage_id/review", "./src/templates/pages/review.hbs")
        .unwrap();

    let app = Router::new()
        .route("/", get(get_index))
        .route("/:passage_id/review", get(get_review))
        // Create the application state
        .with_state(AppState {
            engine: Engine::from(hbs),
            passages: Vec::from([
                Passage {
                    id: 1,
                    reference: String::from("Genesis 1:1-5"),
                    level: 1,
                    // review_date: std::time::Date
                },
                Passage {
                    id: 2,
                    reference: String::from("Genesis 2:5-10"),
                    level: 2,
                    // review_date: std::time::Date
                },
            ]),
        });
    println!("See example: http://127.0.0.1:8080/example");

    let listener = TcpListener::bind((Ipv4Addr::LOCALHOST, 8080))
        .await
        .unwrap();
    serve(listener, app.into_make_service()).await.unwrap();
}
