CREATE TABLE IF NOT EXISTS public.game_settings (
    id SERIAL PRIMARY KEY,
    number_of_words_per_player INTEGER NOT NULL,
    time_per_player INTEGER NOT NULL
);


CREATE TABLE IF NOT EXISTS public.teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    members VARCHAR(255)[] NOT NULL
);



CREATE TABLE IF NOT EXISTS public.game_rounds (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);


INSERT INTO public.game_rounds (name) 
VALUES ('not_started'), ('explain'), ('uno'), ('charades'), ('finished')
ON CONFLICT DO NOTHING;


CREATE TABLE IF NOT EXISTS public.game_state (
    id SERIAL PRIMARY KEY,
    round_id INTEGER NOT NULL REFERENCES public.game_rounds(id) ON DELETE SET NULL,
    player_order VARCHAR(255)[] NOT NULL,
    current_player_index INTEGER NOT NULL

);



CREATE TABLE IF NOT EXISTS public.game (
    id SERIAL PRIMARY KEY,
    game_settings_id INTEGER NOT NULL REFERENCES public.game_settings(id) ON DELETE SET NULL,
    game_state_id INTEGER NOT NULL REFERENCES public.game_state(id) ON DELETE SET NULL,
    teams_id INTEGER[] NOT NULL,
    game_words_state_id INTEGER
);


CREATE TABLE IF NOT EXISTS public.user_words (
    id SERIAL,
    game_id INTEGER NOT NULL REFERENCES public.game(id) ON DELETE SET NULL,
    player_name VARCHAR(255) NOT NULL,
    words VARCHAR(255)[] NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS user_words_unique_idx 
    ON public.user_words (game_id, player_name);


CREATE TABLE IF NOT EXISTS public.rooms (
    id SERIAL PRIMARY KEY,
    room_code VARCHAR(16),
    password VARCHAR(255),
    created_at TIMESTAMP DEFAULT (now() AT TIME ZONE 'UTC'),
    owner_name VARCHAR(255),
    max_members INTEGER,
    current_game_id INTEGER REFERENCES public.game(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS public.room_members (
    room_id INTEGER REFERENCES public.rooms(id) ON DELETE CASCADE,
    player_name VARCHAR(255),
    PRIMARY KEY (room_id, player_name)
);

CREATE TABLE IF NOT EXISTS public.game_words_state (
    id SERIAL PRIMARY KEY,
    words_order VARCHAR(255)[] NOT NULL,
    current_word_index INTEGER NOT NULL
);
