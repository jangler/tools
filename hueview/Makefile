CC = cc
CFLAGS = -Wall -Wextra -Werror -pedantic -I/usr/include/SDL2
LFLAGS = -lSDL2
OBJS = hueview.c
OBJ_NAME = hueview

INSTALL_DIR = /usr/local

all: $(OBJS)
	$(CC) $(OBJS) $(CFLAGS) $(LFLAGS) -o $(OBJ_NAME)

clean:
	rm -f $(OBJ_NAME)

install: $(OBJS)
	install -Dm 755 $(OBJ_NAME) $(INSTALL_DIR)/bin
	install -Dm 644 hueview.1 $(INSTALL_DIR)/man/man1

uninstall:
	rm -f $(INSTALL_DIR)/bin/$(OBJ_NAME) $(INSTALL_DIR)/man/man1/hueview.1
