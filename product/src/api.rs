use actix_web::{web, Responder, HttpResponse, Result};
use crate::models::{NewProduct, UpdateProduct};
use crate::db::establish_connection;
use crate::crud;
use log::{info, error, debug};

pub async fn create_product_handler(item: web::Json<NewProduct>) -> Result<impl Responder> {
    debug!("Received create product request: {:?}", item);
    
    let mut conn = match establish_connection() {
        conn => conn,
    };
    
    info!("Creating product: {}", item.name);
    
    match crud::create_product(&mut conn, item.into_inner()) {
        Ok(product) => {
            info!("Successfully created product with ID: {}", product.id);
            Ok(HttpResponse::Created().json(product))
        },
        Err(e) => {
            error!("Error creating product: {:?}", e);
            Ok(HttpResponse::InternalServerError().json(format!("Error creating product: {}", e)))
        },
    }
}

pub async fn get_all_products_handler() -> Result<impl Responder> {
    debug!("Received get all products request");
    
    let mut conn = match establish_connection() {
        conn => conn,
    };
    
    match crud::get_all_products(&mut conn) {
        Ok(products) => {
            info!("Successfully retrieved {} products", products.len());
            Ok(HttpResponse::Ok().json(products))
        },
        Err(e) => {
            error!("Error fetching products: {:?}", e);
            Ok(HttpResponse::InternalServerError().json(format!("Error fetching products: {}", e)))
        },
    }
}

pub async fn get_product_handler(path: web::Path<i32>) -> Result<impl Responder> {
    let product_id = path.into_inner();
    debug!("Received get product request for ID: {}", product_id);
    
    let mut conn = match establish_connection() {
        conn => conn,
    };
    
    match crud::get_product_by_id(&mut conn, product_id) {
        Ok(product) => {
            info!("Successfully retrieved product with ID: {}", product_id);
            Ok(HttpResponse::Ok().json(product))
        },
        Err(diesel::result::Error::NotFound) => {
            info!("Product with ID {} not found", product_id);
            Ok(HttpResponse::NotFound().json("Product not found"))
        },
        Err(e) => {
            error!("Error fetching product with ID {}: {:?}", product_id, e);
            Ok(HttpResponse::InternalServerError().json(format!("Error fetching product: {}", e)))
        },
    }
}

pub async fn update_product_handler(
    path: web::Path<i32>, 
    item: web::Json<UpdateProduct>
) -> Result<impl Responder> {
    let product_id = path.into_inner();
    debug!("Received update product request for ID: {} with data: {:?}", product_id, item);
    
    let mut conn = match establish_connection() {
        conn => conn,
    };
    
    match crud::update_product(&mut conn, product_id, item.into_inner()) {
        Ok(product) => {
            info!("Successfully updated product with ID: {}", product_id);
            Ok(HttpResponse::Ok().json(product))
        },
        Err(diesel::result::Error::NotFound) => {
            info!("Product with ID {} not found for update", product_id);
            Ok(HttpResponse::NotFound().json("Product not found"))
        },
        Err(e) => {
            error!("Error updating product with ID {}: {:?}", product_id, e);
            Ok(HttpResponse::InternalServerError().json(format!("Error updating product: {}", e)))
        },
    }
}

pub async fn delete_product_handler(path: web::Path<i32>) -> Result<impl Responder> {
    let product_id = path.into_inner();
    debug!("Received delete product request for ID: {}", product_id);
    
    let mut conn = match establish_connection() {
        conn => conn,
    };
    
    match crud::delete_product(&mut conn, product_id) {
        Ok(0) => {
            info!("Product with ID {} not found for deletion", product_id);
            Ok(HttpResponse::NotFound().json("Product not found"))
        },
        Ok(_) => {
            info!("Successfully deleted product with ID: {}", product_id);
            Ok(HttpResponse::Ok().json("Product deleted successfully"))
        },
        Err(e) => {
            error!("Error deleting product with ID {}: {:?}", product_id, e);
            Ok(HttpResponse::InternalServerError().json(format!("Error deleting product: {}", e)))
        },
    }
}
