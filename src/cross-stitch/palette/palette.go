package palette

import (
  "encoding/csv"
  "os"
  "strconv"
)

type Thread struct {
  ID int
  Name string
  R uint8
  G uint8
  B uint8
}

func GreyPalette () ([]Thread, error) {
  return palette("../../palette/black-white-grey.csv")
}

func DMCPalette () ([]Thread, error) {
  return palette("../../palette/dmc-floss.csv")
}

func palette (path string) ([]Thread, error) {
  file, err := os.Open(path)
  if err != nil { return nil, err }
  defer file.Close()

  // convert dmc data to csv hash
  reader := csv.NewReader(file)
  reader.Comma = ','

  // Record: [["Floss#","Description","Red","Green","Blue"],...]
  record, err := reader.ReadAll()
  if err != nil {
    return nil, err
  }

  t := []Thread{}
  for i := 0; i < len(record); i++ {
    id, _ := strconv.Atoi(record[i][0])
    tr, _ := strconv.Atoi(record[i][2])
    tg, _ := strconv.Atoi(record[i][3])
    tb, _ := strconv.Atoi(record[i][4])
    thread := Thread{id, record[i][1], uint8(tr), uint8(tg), uint8(tb)}

    t = append(t, thread)
  }

  return t, nil
}
