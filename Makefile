##
## backend
## File description:
## Makefile
##

GO	=	go

NAME	=	app

SRCDIR	=	init

SRC		=	main.go

SRC			:= $(addprefix $(SRCDIR)/, $(SRC))

GOFLAGS =	--trimpath --mod=vendor

all: $(NAME)

$(NAME):
	$(GO) mod vendor
	$(GO) build $(GOFLAGS) -o $(NAME) $(SRC)

fclean:
	rm -f ./$(NAME)

re:	fclean all