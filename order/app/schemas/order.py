from pydantic import BaseModel, Field, ConfigDict
from typing import List, Optional, Dict, Any
from datetime import datetime
from decimal import Decimal
from uuid import UUID
from enum import Enum

class OrderStatus(str, Enum):
    PENDING = "pending"
    CONFIRMED = "confirmed"
    PROCESSING = "processing"
    SHIPPED = "shipped"
    DELIVERED = "delivered"
    CANCELLED = "cancelled"
    REFUNDED = "refunded"

class AddressBase(BaseModel):
    street: str
    city: str
    state: str
    postal_code: str
    country: str
    full_name: str
    phone: Optional[str] = None

class OrderItemCreate(BaseModel):
    product_id: int
    quantity: int = Field(gt=0)

class OrderItemResponse(BaseModel):
    id: UUID
    product_id: int
    quantity: int
    unit_price: Decimal
    total_price: Decimal
    product_name: Optional[str] = None
    product_description: Optional[str] = None
    created_at: datetime

class OrderCreate(BaseModel):
    items: List[OrderItemCreate]
    shipping_address: AddressBase
    billing_address: Optional[AddressBase] = None
    notes: Optional[str] = None

class OrderResponse(BaseModel):
    id: UUID
    user_id: UUID
    status: str
    total_amount: Decimal
    currency: str
    shipping_address: Dict[str, Any]
    billing_address: Optional[Dict[str, Any]]
    notes: Optional[str]
    created_at: datetime
    updated_at: datetime
    items: List[OrderItemResponse]

class OrderListResponse(BaseModel):
    id: UUID
    status: str
    total_amount: Decimal
    currency: str
    created_at: datetime
    items_count: int

class OrderStatusUpdate(BaseModel):
    status: OrderStatus
    notes: Optional[str] = None
