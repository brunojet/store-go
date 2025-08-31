# Estrutura de Armazenamento de Imagens no S3

## Diretório Raiz
```
app_images/
```

## Organização por Aplicativo
Cada aplicativo possui um diretório próprio identificado por seu nome ou slug:
```
app_images/aplicativo/{nome_aplicativo}/
```

## Ícones
Os ícones são armazenados em subdiretórios conforme o tipo e tamanho:
```
app_images/aplicativo/{nome_aplicativo}/icons/{hash_imagem}/{thumb|pequeno|medio|grande|original}.{extensao}
```

## Screenshots
Screenshots são organizados por contexto e tamanho:
```
app_images/aplicativo/{nome_aplicativo}/screenshots/{hash_imagem}/{thumb|pequeno|medio|grande|original}.{extensao}
```

## Banners
Banners são armazenados conforme o tipo e proporção:
```
app_images/aplicativo/{nome_aplicativo}/banners/{hash_imagem}/{thumb|pequeno|medio|grande|original}.{extensao}
```

## Observações

## Exemplos de Estrutura

### Ícones
```
app_images/aplicativo/loja_super/icons/abc123def456/thumb.png        # Miniatura (ex: 64x64 px)
app_images/aplicativo/loja_super/icons/abc123def456/pequeno.png      # Pequeno (ex: 128x128 px)
app_images/aplicativo/loja_super/icons/abc123def456/medio.png        # Médio (ex: 192x192 px)
app_images/aplicativo/loja_super/icons/abc123def456/grande.png       # Grande (ex: 512x512 px)
app_images/aplicativo/loja_super/icons/abc123def456/original.png     # Original, sem redimensionamento
```

### Screenshots
```
app_images/aplicativo/loja_super/screenshots/xyz789abc012/thumb.jpg      # Miniatura (ex: 320x180 px)
app_images/aplicativo/loja_super/screenshots/xyz789abc012/medio.jpg      # Médio (ex: 1080x1920 px)
app_images/aplicativo/loja_super/screenshots/xyz789abc012/grande.jpg     # Grande (ex: 1920x1080 px)
app_images/aplicativo/loja_super/screenshots/xyz789abc012/original.jpg   # Original, sem redimensionamento
```

### Banners
```
app_images/aplicativo/loja_super/banners/qwe456rty789/thumb.png      # Miniatura (ex: 300x100 px)
app_images/aplicativo/loja_super/banners/qwe456rty789/medio.png      # Médio (ex: 1200x300 px)
app_images/aplicativo/loja_super/banners/qwe456rty789/grande.png     # Grande (ex: 1920x480 px)
app_images/aplicativo/loja_super/banners/qwe456rty789/original.png   # Original, sem redimensionamento
```
