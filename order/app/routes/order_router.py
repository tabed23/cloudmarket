from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session
from typing import List
from uuid import UUID
from app.schemas.order import OrderCreate, OrderResponse, OrderListResponse, OrderStatusUpdate, OrderStatusHistoryResponse
from app.services.order_svc import OrderService
from app.database import get_db

router = APIRouter()

# Mock authentication - replace with actual auth service integration
def get_current_user():
    # This should validate JWT token and return user info
    return UUID("123e4567-e89b-12d3-a456-426614174000")

@router.post("/", response_model=OrderResponse)
async def create_order(
    order_data: OrderCreate,
    current_user: UUID = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Create a new order"""
    service = OrderService(db)
    try:
        order = await service.create_order(current_user, order_data)
        return order
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=500, detail="Internal server error")

@router.get("/", response_model=List[OrderListResponse])
async def get_user_orders(
    skip: int = Query(0, ge=0),
    limit: int = Query(100, ge=1, le=100),
    current_user: UUID = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get current user's orders"""
    service = OrderService(db)
    orders = service.get_user_orders(current_user, skip, limit)
    
    return [
        OrderListResponse(
            id=order.id,
            status=order.status,
            total_amount=order.total_amount,
            currency=order.currency,
            created_at=order.created_at,
            items_count=len(order.items)
        )
        for order in orders
    ]

@router.get("/{order_id}", response_model=OrderResponse)
async def get_order(
    order_id: UUID,
    current_user: UUID = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get a specific order"""
    service = OrderService(db)
    order = service.get_order_by_id(order_id, current_user)
    
    if not order:
        raise HTTPException(status_code=404, detail="Order not found")
    
    return order

@router.put("/{order_id}/status", response_model=OrderResponse)
async def update_order_status(
    order_id: UUID,
    status_update: OrderStatusUpdate,
    current_user: UUID = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Update order status"""
    service = OrderService(db)
    order = service.update_order_status(order_id, status_update, current_user)
    
    if not order:
        raise HTTPException(status_code=404, detail="Order not found")
    
    return order

@router.get("/{order_id}/history", response_model=List[OrderStatusHistoryResponse])
async def get_order_status_history(
    order_id: UUID,
    current_user: UUID = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """Get order status history"""
    service = OrderService(db)
    
    # Verify order belongs to user
    order = service.get_order_by_id(order_id, current_user)
    if not order:
        raise HTTPException(status_code=404, detail="Order not found")
    
    history = service.get_order_status_history(order_id)
    return history


