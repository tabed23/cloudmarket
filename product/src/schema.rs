diesel::table! {
    products (id) {
        id -> Int4,
        name -> Varchar,
        price -> Float8,
        description -> Nullable<Text>,
        created_at -> Timestamp,
        updated_at -> Timestamp,
    }
}