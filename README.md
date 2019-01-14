# unmagic

Unmagic can reveal the file type based on magic numbers searched in head and tail of file. It uses `signatures.json` as a database which is formatted this way :


```
[
  {
    "header": "ffd8ffe0xxxx4a46494600",
    "trailer": "ffd9",
    "description": "JFIF, JPE, JPEG, JPG, JPEG/JFIF graphics file"
  },
  {
    "header": "89504E470D0A1A0A",
    "trailer": "49454E44AE426082",
    "description": "PNG : Portable Network Graphics file"
  }
]
```

