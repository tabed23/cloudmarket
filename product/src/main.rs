use actix_web::{middleware::Logger, web, App, HttpServer};
use env_logger::Env;
use product::{api, db::establish_connection};
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::Builder::from_env(Env::default().default_filter_or("debug")).init();
    // Test database connection
    match establish_connection() {
        conn => {
            println!("Successfully connected to the database.");
            drop(conn);
        }
    }

    println!("Starting server on http://127.0.0.1:8080");

    HttpServer::new(|| {
        App::new()
            .wrap(Logger::default())
            .route("/products", web::post().to(api::create_product_handler))
            .route("/products", web::get().to(api::get_all_products_handler))
            .route("/products/{id}", web::get().to(api::get_product_handler))
            .route("/products/{id}", web::put().to(api::update_product_handler))
            .route(
                "/products/{id}",
                web::delete().to(api::delete_product_handler),
            )
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
