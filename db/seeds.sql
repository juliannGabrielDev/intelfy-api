INSERT INTO genres (id, name) VALUES 
('genre_1234567890abcde', 'Tele-Rock'),
('genre_qwertyuiop12345', 'Lo-Fi Engineering');

INSERT INTO users (id, username, email, password_hash) VALUES 
('user_admin_000000001', 'tele_user', 'user@tele.edu.co', 'hash_de_password_seguro'),
('art_88888888888888888', 'artist_user', 'artist@tele.edu.co', 'hash_de_password_seguro');

INSERT INTO artists (id, name, bio, cover_url) VALUES 
('art_88888888888888888', 'The Tele-Informatics Duo', 'Banda formada en los laboratorios de la facultad.', 'artists/covers/duo.jpg');

INSERT INTO albums (id, name, artist_id, cover_url) VALUES 
('alb_99999999999999999', 'First Protocol', 'art_88888888888888888', 'albums/covers/first_protocol.jpg');

INSERT INTO songs (id, name, duration_seconds, audio_url, cover_url, album_id, genre_id) VALUES 
('song_0000000000000001', 'TCP Handshake', 180, 'songs/audio/tcp.mp3', 'songs/covers/tcp.jpg', 'alb_99999999999999999', 'genre_1234567890abcde'),
('song_0000000000000002', 'UDP Stream', 210, 'songs/audio/udp.mp3', 'songs/covers/udp.jpg', 'alb_99999999999999999', 'genre_qwertyuiop12345');