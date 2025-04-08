containers=(arb-nginx arb-backend-1 arb-backend-2)

for container in "${containers[@]}"; do
  
    echo "==========================================================================================="
    echo "$container"
    echo "==========================================================================================="
    docker logs "$container"
    
    echo
done