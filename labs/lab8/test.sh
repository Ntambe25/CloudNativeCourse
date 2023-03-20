
list()
{
    echo "List of items in Database: "
    curl -w "\n" "http://localhost:8000/list"
}

price()
{
    echo "Price of item: pant"
    curl -w "\n" "http://localhost:8000/price?item=pant"
}

create()
{
    echo "adding item: $1 in database of price: $2"
    curl -w "\n" "http://localhost:8000/create?item=$1&price=$2"
}

update()
{
    echo "updating price of item: $1 with price: $2"
    curl -w "\n" "http://localhost:8000/update?item=$1&price=$2"
}

delete()
{
    echo "deleting item: $1 from database"
    curl -w "\n" "http://localhost:8000/delete?item=$1"
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