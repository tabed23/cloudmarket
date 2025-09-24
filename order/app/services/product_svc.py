import httpx
from typing import Optional, Dict, Any
from app.core.config import settings

class ProductService:
    def __init__(self):
        self.base_url = settings.product_svc_url
        
    async def get_product(self, product_id: int) -> Optional[Dict[str, Any]]:
        """Get product details from the product service"""
        async with httpx.AsyncClient() as client:
            try:
                response = await client.get(f"{self.base_url}/products/{product_id}")
                if response.status_code == 200:
                    return response.json()
                return None
            except httpx.RequestError:
                return None
    
    async def get_products(self, product_ids: list[int]) -> Dict[int, Dict[str, Any]]:
        """Get multiple products from the product service"""
        products = {}
        async with httpx.AsyncClient() as client:
            for product_id in product_ids:
                try:
                    response = await client.get(f"{self.base_url}/products/{product_id}")
                    if response.status_code == 200:
                        products[product_id] = response.json()
                except httpx.RequestError:
                    continue
        return products
