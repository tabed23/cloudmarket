use chrono::NaiveDateTime;
use serde::{Deserialize, Serialize};
use diesel::{Queryable, Insertable, AsChangeset};
use crate::schema::products;

#[derive(Queryable, Serialize, Deserialize, Debug, Clone)]
pub struct Product {
    pub id: i32,
    pub name: String,
    pub price: f64,
    pub description: Option<String>,
    pub created_at: NaiveDateTime,
    pub updated_at: NaiveDateTime,
}

#[derive(Insertable, Serialize, Deserialize, Debug)]
#[diesel(table_name = products)]
pub struct NewProduct {
    pub name: String,
    pub price: f64,
    pub description: Option<String>,
}

#[derive(AsChangeset, Serialize, Deserialize, Debug)]
#[diesel(table_name = products)]
pub struct UpdateProduct {
    pub name: Option<String>,
    pub price: Option<f64>,
    pub description: Option<String>,
}