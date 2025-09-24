use crate::models::{Product, NewProduct, UpdateProduct};
use diesel::prelude::*;
use chrono::Utc;
use log::{debug, error};

pub fn create_product(conn: &mut PgConnection, new_product: NewProduct) -> Result<Product, diesel::result::Error> {
    use crate::schema::products::dsl::*;
    
    debug!("Inserting new product into database: {:?}", new_product);
    
    let result = diesel::insert_into(products)
        .values(&new_product)
        .get_result(conn);
        
    match &result {
        Ok(product) => debug!("Successfully inserted product: {:?}", product),
        Err(e) => error!("Database error during product insertion: {:?}", e),
    }
    
    result
}

pub fn get_all_products(conn: &mut PgConnection) -> Result<Vec<Product>, diesel::result::Error> {
    use crate::schema::products::dsl::*;
    
    debug!("Fetching all products from database");
    
    let result = products.load::<Product>(conn);
    
    match &result {
        Ok(products_vec) => debug!("Successfully fetched {} products", products_vec.len()),
        Err(e) => error!("Database error during products fetch: {:?}", e),
    }
    
    result
}

pub fn get_product_by_id(conn: &mut PgConnection, product_id: i32) -> Result<Product, diesel::result::Error> {
    use crate::schema::products::dsl::*;
    
    debug!("Fetching product with ID: {}", product_id);
    
    let result = products
        .filter(id.eq(product_id))
        .first::<Product>(conn);
        
    match &result {
        Ok(product) => debug!("Successfully fetched product: {:?}", product),
        Err(diesel::result::Error::NotFound) => debug!("Product with ID {} not found", product_id),
        Err(e) => error!("Database error during product fetch: {:?}", e),
    }
    
    result
}

pub fn update_product(conn: &mut PgConnection, product_id: i32, update_data: UpdateProduct) -> Result<Product, diesel::result::Error> {
    use crate::schema::products::dsl::*;
    
    debug!("Updating product with ID: {} with data: {:?}", product_id, update_data);
    
    let result = diesel::update(products.filter(id.eq(product_id)))
        .set((
            &update_data,
            updated_at.eq(Utc::now().naive_utc())
        ))
        .get_result(conn);
        
    match &result {
        Ok(product) => debug!("Successfully updated product: {:?}", product),
        Err(diesel::result::Error::NotFound) => debug!("Product with ID {} not found for update", product_id),
        Err(e) => error!("Database error during product update: {:?}", e),
    }
    
    result
}

pub fn delete_product(conn: &mut PgConnection, product_id: i32) -> Result<usize, diesel::result::Error> {
    use crate::schema::products::dsl::*;
    
    debug!("Deleting product with ID: {}", product_id);
    
    let result = diesel::delete(products.filter(id.eq(product_id)))
        .execute(conn);
        
    match &result {
        Ok(count) => debug!("Successfully deleted {} product(s)", count),
        Err(e) => error!("Database error during product deletion: {:?}", e),
    }
    
    result
}