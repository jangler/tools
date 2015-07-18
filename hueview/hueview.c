#include <SDL.h>
#include <stdio.h>
#include <string.h>

char *title = "hueview";
int width = 320, height = 240;

SDL_Window *window = NULL;
SDL_Renderer *renderer = NULL;

int colorc = 0;
SDL_Color *colorv = NULL;

/* returns 1 on success or 0 on failure */
int init(void)
{
	if (SDL_Init(SDL_INIT_VIDEO | SDL_INIT_AUDIO) < 0) {
		fprintf(stderr, "SDL_Init: %s\n", SDL_GetError());
		return 0;
	}
	
	window = SDL_CreateWindow(title, SDL_WINDOWPOS_UNDEFINED,
			SDL_WINDOWPOS_UNDEFINED, width, height,
			SDL_WINDOW_RESIZABLE);
	if (!window) {
		fprintf(stderr, "SDL_CreateWindow: %s\n", SDL_GetError());
		return 0;
	}

	renderer = SDL_CreateRenderer(window, -1, 0);
	if (!renderer) {
		fprintf(stderr, "SDL_CreateRenderer: %s\n", SDL_GetError());
		return 0;
	}

	return 1;
}

void draw(void)
{
	int i;
	SDL_Rect rect;
	
	for (i = 0; i < colorc; ++i) {
		rect = (SDL_Rect){ i * width / colorc, 0, width / colorc,
				height };
		SDL_SetRenderDrawColor(renderer, colorv[i].r, colorv[i].g,
				colorv[i].b, colorv[i].a);
		SDL_RenderFillRect(renderer, &rect);
	}
	SDL_RenderPresent(renderer);
}

void eventloop(void)
{
	SDL_Event event;
	
	for (;;) {
		if (!SDL_WaitEvent(&event)) {
			fprintf(stderr, "SDL_WaitEvent: %s\n", SDL_GetError());
			break;
		}
		
		if (event.type == SDL_QUIT)
			break;
		else if (event.type == SDL_WINDOWEVENT) {
			if (event.window.event == SDL_WINDOWEVENT_RESIZED) {
				width = event.window.data1;
				height = event.window.data2;
			}
			draw();
		}
	}
}

int main(int argc, char *argv[])
{
	int i, v;
	
	if (argc < 2) {
		fprintf(stderr, "usage: %s color...\n", argv[0]);
		return 1;
	}
	
	/* init colors */
	colorc = argc - 1;
	colorv = malloc(sizeof(SDL_Color) * colorc);
	if (!colorv) {
		perror("could not init colors");
		return 1;
	}
	for (i = 1; i < argc; ++i) {
		v = strtol(argv[i], NULL, 16);
		colorv[i - 1] = (SDL_Color){
			.r = (v >> 16) & 0xff,
			.g = (v >> 8) & 0xff,
			.b = v & 0xff,
			.a = 0xff
		};
	}
	
	if (init()) {
		draw();
		eventloop();
	}
	
	if (renderer)
		SDL_DestroyRenderer(renderer);
	if (window)
		SDL_DestroyWindow(window);
	if (colorv)
		free(colorv);
	SDL_Quit();

	return 0;
}
