from pydantic_settings import BaseSettings
from typing import Optional

class Settings(BaseSettings):
   database_url: str
   secret_key: str
   algorithm: str
   access_token_expire_minutes: int
   
   product_svc_url : str
   auth_svc_url : str


class Config:
    env_file = ".env"
    
settings = Settings()

print("Database URL:", settings.database_url)
print("Product Service URL:", settings.product_svc_url)
print("Auth Service URL:", settings.auth_svc_url)