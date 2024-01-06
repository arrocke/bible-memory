use std::net::Ipv4Addr;

use axum::{
    extract::{FromRef, Path},
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

#[derive(Debug, Serialize)]
pub struct Person {
    name: String,
}

async fn get_name(engine: AppEngine, Key(key): Key, Path(name): Path<String>) -> impl IntoResponse {
    let person = Person { name };
    RenderHtml(key, engine, person)
}

async fn get_index(engine: AppEngine) -> impl IntoResponse {
    let person = Person {
        name: String::from("Index"),
    };
    RenderHtml(Key(String::from("/:name")), engine, person)
}

#[derive(Clone, FromRef)]
struct AppState {
    engine: AppEngine,
}

#[tokio::main]
async fn main() {
    let mut hbs = Handlebars::new();
    hbs.register_template_file("/:name", "./src/templates/index.hbs")
        .unwrap();

    let app = Router::new()
        .route("/:name", get(get_name))
        .route("/", get(get_index))
        // Create the application state
        .with_state(AppState {
            engine: Engine::from(hbs),
        });
    println!("See example: http://127.0.0.1:8080/example");

    let listener = TcpListener::bind((Ipv4Addr::LOCALHOST, 8080))
        .await
        .unwrap();
    serve(listener, app.into_make_service()).await.unwrap();
}
