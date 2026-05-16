INSERT INTO genres (id, name) VALUES 
('genre_1234567890abcde', 'Tele-Rock'),
('genre_qwertyuiop12345', 'Lo-Fi Engineering');

INSERT INTO artists (id, name, bio) VALUES 
('art_88888888888888888', 'The Tele-Informatics Duo', 'Banda formada en los laboratorios de la facultad.');

INSERT INTO users (id, username, email, password_hash) VALUES 
('user_admin_000000001', 'tele_user', 'user@tele.edu.co', 'hash_de_password_seguro');

INSERT INTO albums (id, name, artist_id) VALUES 
('alb_99999999999999999', 'First Protocol', 'art_88888888888888888');

INSERT INTO songs (id, name, duration_seconds, audio_url, album_id, genre_id) VALUES 
('song_0000000000000001', 'TCP Handshake', 180, 'https://cdn.intelfy.com/s1.mp3', 'alb_99999999999999999', 'genre_1234567890abcde'),
('song_0000000000000002', 'UDP Stream', 210, 'https://cdn.intelfy.com/s2.mp3', 'alb_99999999999999999', 'genre_qwertyuiop12345');