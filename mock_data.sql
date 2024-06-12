BEGIN;

-- Create 100 articles with random data
INSERT INTO article (articleid, title, text)
SELECT
  GENERATE_SERIES(1, 100) AS articleid,
  MD5(RANDOM()::TEXT) AS title,
  MD5(RANDOM()::TEXT) AS text;

-- Create 100 comments for each article with random data
INSERT INTO comment (commentid, text, rating, articleid)
SELECT
  GENERATE_SERIES(1, 10000) AS commentid,
  MD5(RANDOM()::TEXT) AS text,
  FLOOR(RANDOM() * 100 + 1) AS rating, -- Generate random rating between 1 and 100
  FLOOR(RANDOM() * 100 + 1) AS articleid;

COMMIT;
