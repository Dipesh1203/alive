import json
import uuid
import random
from datetime import datetime, timedelta

def generate_mock_aws_logs(total_requests=5):
    services = ["api-gateway", "auth-service", "payment-service", "inventory-service"]
    exceptions = [
        {"type": "java.sql.SQLTransientConnectionException", "msg": "Connection is not available, request timed out after 3000ms."},
        {"type": "redis.clients.jedis.exceptions.JedisConnectionException", "msg": "Could not get a resource from the pool"},
        {"type": "com.amazonaws.services.s3.model.AmazonS3Exception", "msg": "Access Denied (Service: Amazon S3; Status Code: 403)"}
    ]
    
    logs = []
    base_time = datetime.utcnow()

    for i in range(total_requests):
        # Create matching correlation keys for the entire lifecycle of this single request
        request_id = str(uuid.uuid4())
        trace_id = f"1-{hex(random.getrandbits(32))[2:]}-{hex(random.getrandbits(64))[2:]}"
        user_id = f"usr_{random.randint(1000, 9999)}"
        order_id = f"ord_{random.randint(50000, 90000)}"
        
        # Decide if this specific request lifecycle is going to fail
        is_failure_chain = random.choice([True, False])
        
        # Step 1: Request Hits Gateway
        base_time += timedelta(milliseconds=random.randint(10, 50))
        logs.append({
            "timestamp": base_time.isoformat() + "Z",
            "aws_request_id": request_id,
            "log_level": "INFO",
            "service": "api-gateway",
            "environment": "production",
            "trace_id": trace_id,
            "message": f"Incoming POST request to /v1/checkout",
            "context": {"user_id": user_id}
        })
        
        # Step 2: Request passes downstream to Payment Service
        base_time += timedelta(milliseconds=random.randint(20, 100))
        logs.append({
            "timestamp": base_time.isoformat() + "Z",
            "aws_request_id": request_id,
            "log_level": "INFO",
            "service": "payment-service",
            "environment": "production",
            "trace_id": trace_id,
            "message": f"Processing transaction for checkout asset: {order_id}",
            "context": {"amount": round(random.uniform(10.0, 500.0), 2), "currency": "INR"}
        })

        if is_failure_chain:
            # Step 3a: Issue an early warning sign
            base_time += timedelta(seconds=random.randint(1, 2))
            logs.append({
                "timestamp": base_time.isoformat() + "Z",
                "aws_request_id": request_id,
                "log_level": "WARN",
                "service": "payment-service",
                "environment": "production",
                "trace_id": trace_id,
                "message": "Resource constraint detected. Upstream database queue length exceeds threshold.",
                "context": {"current_queue_depth": random.randint(85, 120)}
            })
            
            # Step 4a: Drop the critical application exception block
            selected_err = random.choice(exceptions)
            base_time += timedelta(seconds=1)
            logs.append({
                "timestamp": base_time.isoformat() + "Z",
                "aws_request_id": request_id,
                "log_level": "ERROR",
                "service": "payment-service",
                "environment": "production",
                "trace_id": trace_id,
                "message": "Critical dependency failure inside runtime loop processing.",
                "error": {
                    "exception": selected_err["type"],
                    "message": selected_err["msg"],
                    "stack_trace": f"at com.app.internal.Executor.run(Executor.java:14)\n at com.app.payment.Handler.process(Handler.java:89)"
                }
            })
            
            # Step 5a: Gateway logs the final response failure to client
            base_time += timedelta(milliseconds=15)
            logs.append({
                "timestamp": base_time.isoformat() + "Z",
                "aws_request_id": request_id,
                "log_level": "FATAL",
                "service": "api-gateway",
                "environment": "production",
                "trace_id": trace_id,
                "message": "Upstream microservice returned response status code 500.",
                "context": {"http_status": 500, "duration_ms": random.randint(3000, 4500)}
            })
        else:
            # Step 3b: Clean execution path completes
            base_time += timedelta(milliseconds=random.randint(50, 200))
            logs.append({
                "timestamp": base_time.isoformat() + "Z",
                "aws_request_id": request_id,
                "log_level": "INFO",
                "service": "payment-service",
                "environment": "production",
                "trace_id": trace_id,
                "message": f"Successfully completed execution flow for order {order_id}",
                "context": {"status": "SUCCESS"}
            })

    # Sort all generated logs globally by timestamp to mimic a real unified log stream
    logs.sort(key=lambda x: x["timestamp"])
    return logs

# Execute and write to file
mock_dataset = generate_mock_aws_logs(total_requests=50)

with open("mock_aws_production_logs.json", "w") as f:
    json.dump(mock_dataset, f, indent=2)

print(f"Generated {len(mock_dataset)} sequential log entries successfully.")