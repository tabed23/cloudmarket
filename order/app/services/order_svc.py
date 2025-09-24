from sqlalchemy.orm import Session
from sqlalchemy import and_, desc
from typing import List, Optional
from uuid import UUID
from decimal import Decimal
from app.models.order import Order, OrderItem, OrderStatusHistory
from app.schemas.order import (
    OrderCreate, OrderUpdate, OrderStatusUpdate, 
    OrderCalculation, OrderCalculationResponse
)
from app.services.product_svc import product_service

class OrderService:
    def __init__(self, db: Session):
        self.db = db
    
    async def create_order(self, user_id: UUID, order_data: OrderCreate) -> Order:
        """Create a new order"""
        # Get product details
        product_ids = [item.product_id for item in order_data.items]
        products = await product_service.get_products(product_ids)
        
        # Calculate order total
        total_amount = Decimal('0')
        order_items = []
        
        for item_data in order_data.items:
            product = products.get(item_data.product_id)
            if not product:
                raise ValueError(f"Product {item_data.product_id} not found")
            
            unit_price = Decimal(str(product['price']))
            total_price = unit_price * item_data.quantity
            total_amount += total_price
            
            order_items.append({
                'product_id': item_data.product_id,
                'quantity': item_data.quantity,
                'unit_price': unit_price,
                'total_price': total_price,
                'product_name': product['name'],
                'product_description': product.get('description')
            })
        
        # Create order
        order = Order(
            user_id=user_id,
            total_amount=total_amount,
            shipping_address=order_data.shipping_address.model_dump(),
            billing_address=order_data.billing_address.model_dump() if order_data.billing_address else None,
            notes=order_data.notes
        )
        
        self.db.add(order)
        self.db.flush()  # Get the order ID
        
        # Create order items
        for item_data in order_items:
            order_item = OrderItem(order_id=order.id, **item_data)
            self.db.add(order_item)
        
        # Create initial status history
        status_history = OrderStatusHistory(
            order_id=order.id,
            status="pending",
            notes="Order created"
        )
        self.db.add(status_history)
        
        self.db.commit()
        self.db.refresh(order)
        return order
    
    def get_user_orders(self, user_id: UUID, skip: int = 0, limit: int = 100) -> List[Order]:
        """Get orders for a specific user"""
        return self.db.query(Order).filter(
            Order.user_id == user_id
        ).order_by(desc(Order.created_at)).offset(skip).limit(limit).all()
    
    def get_order_by_id(self, order_id: UUID, user_id: Optional[UUID] = None) -> Optional[Order]:
        """Get order by ID, optionally filtered by user"""
        query = self.db.query(Order).filter(Order.id == order_id)
        if user_id:
            query = query.filter(Order.user_id == user_id)
        return query.first()
    
    def update_order_status(self, order_id: UUID, status_update: OrderStatusUpdate, changed_by: Optional[UUID] = None) -> Optional[Order]:
        """Update order status"""
        order = self.db.query(Order).filter(Order.id == order_id).first()
        if not order:
            return None
        
        # Update order status
        old_status = order.status
        order.status = status_update.status
        
        # Add status history
        status_history = OrderStatusHistory(
            order_id=order.id,
            status=status_update.status,
            notes=status_update.notes or f"Status changed from {old_status} to {status_update.status}",
            changed_by=changed_by
        )
        self.db.add(status_history)
        
        self.db.commit()
        self.db.refresh(order)
        return order
