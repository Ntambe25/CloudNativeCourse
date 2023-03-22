
list()
{
    echo "List of items in Database: "
    curl -w "\n" "http://localhost:8000/list"
}

create()
{
    echo "adding item: socks in database of price: $10"
    curl -w "\n" "http://localhost:8000/create?item=socks&price=10"
}

price()
{
    echo "Price of item: socks"
    curl -w "\n" "http://localhost:8000/price?item=socks"
}


update()
{
    echo "updating price of item: socks with price: $25"
    curl -w "\n" "http://localhost:8000/update?item=socks&price=25"
}

delete()
{
    echo "deleting item: socks from database"
    curl -w "\n" "http://localhost:8000/delete?item=socks"
}
 
printf "Enter input--->\nlist\nprice\ncreate\nupdate\ndelete\n"
read OPTION
case $OPTION in
  list)
    list
    ;;

  price)
    read -p "Enter item: " ITEM
    price $ITEM
    ;;

  create)
    read -p "Enter item: " ITEM
    echo "Enter price of item: "
    read PRICE
    create $ITEM $PRICE
    ;;

  update)
    read -p "Enter item: " ITEM
    echo "Enter price of item: "
    read PRICE
    update $ITEM $PRICE
    ;;

    delete)
    read -p "Enter item: " ITEM
    delete $ITEM
    ;;

  *)
    echo -n "unknown option: $OPTION"
    ;;
esac