pub mod get_new_passage;
pub mod get_passage_review;
pub mod get_passages;
pub mod post_passage;

use axum::extract::FromRef;
use sqlx::{postgres::Postgres, Pool};

pub type DbPool = Pool<Postgres>;

#[derive(Clone, FromRef)]
pub struct AppState {
    pub db_pool: DbPool,
}
