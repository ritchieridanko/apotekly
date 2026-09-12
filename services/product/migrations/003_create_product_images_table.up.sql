CREATE TABLE product_images (
  id UUID PRIMARY KEY,
  product_id UUID NOT NULL,

  image VARCHAR NOT NULL,
  alt_text VARCHAR,
  sort_order SMALLINT NOT NULL DEFAULT 0,
  
  UNIQUE (product_id, sort_order),
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

-- Index records by product
CREATE INDEX idx_product_images_product ON product_images (product_id);
