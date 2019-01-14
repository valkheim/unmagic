# unmagic

Unmagic reveals the file type based on magic numbers.

<p align="center"> 
  <img src="https://raw.githubusercontent.com/valkheim/unmagic/readme/preview.gif"/>
</p>

It uses `signatures.json` as a database which is formatted this way :

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
