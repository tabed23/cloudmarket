from fastapi import FastAPI, Depends, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager
import uvicorn
from app.database import engine, Base
from app.routes import orders, health
from app.core.config import settings

# Create tables on startup
@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    print("Creating database tables...")
    Base.metadata.create_all(bind=engine)
    print("Database tables created successfully!")
    yield
    # Shutdown
    print("Shutting down...")

# FastAPI app instance
app = FastAPI(
    title="Order Service API",
    description="Microservice for managing orders in an e-commerce platform",
    version="1.0.0",
    lifespan=lifespan
)

# CORS middleware to allow cross-origin requests
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Adjust the allowed origins in production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(orders.router, prefix="/api/v1/orders", tags=["Orders"])

@app.get("/")
async def root():
    return {"message": "Order Service API", "version": "1.0.0"}

# Run the app with Uvicorn (usually handled outside the app code)
if __name__ == "__main__":
    uvicorn.run("app.main:app", host="0.0.0.0", port=8001, reload=True)
