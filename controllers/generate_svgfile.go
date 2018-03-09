package controllers

import (
	"fmt"
	"panda/models"

	"math/rand"
	"time"

	"bufio"
	"os"

	"github.com/astaxie/beego"
)

type GeneratesvgfileController struct {
	beego.Controller
}

/*get random number*/
func generate_rand(number int64) (nRand int64) {
	if number <= 0 {
		return -1
	}

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	return int64(seed.Intn(int(number)))
}

func (c *GeneratesvgfileController) HandlerGenerate() {

	path := c.Generate_svg(1, "/root/gocode/src/panda/svgfile/", "1")

	c.Ctx.WriteString(c.Ctx.Request.Host)

	c.Ctx.WriteString("http://47.92.67.93:8080/svg/" + path)
}

/*generate svg file*/
func (c *GeneratesvgfileController) Generate_svg(flag int, basePath string, petID string) (svgPath string) {

	//svg head
	svg := "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 800 800\">"

	//random for panda base color
	color := generate_rand(7)

	//memset
	var selectflag_times = make([]int, 4, 10)
	var selectflag_check = make([]int, 4, 10)

	for i, _ := range selectflag_check {
		if i > 1 {
			count, _ := models.GetCountBySelectId(int64(i))
			selectflag_check[i] = int(generate_rand(count))
		}
	}

	//query structController
	query := make(map[string]string, 0)

	//Get svg_catatory order by rank
	ml, err := models.GetAllSvgcata(query, []string{}, []string{"rank"}, []string{"asc"}, 0, 20)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range ml {
		catagory_id := v.(models.Svg_catagory).Id
		select_flag := v.(models.Svg_catagory).Select_flag
		bodycolor_flag := v.(models.Svg_catagory).Bodycolor_flag
		percent := v.(models.Svg_catagory).Probability

		switch select_flag {
		case 0:
			//use random color
			if bodycolor_flag == 1 {

				rand := generate_rand(models.GetCountByCatagoryId(catagory_id) / 7)
				if rand == -1 {
					break
				}

				query := make(map[string]string, 0)
				query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
				query["base_color"] = fmt.Sprintf("%v", color)
				query["p_id"] = fmt.Sprintf("%v", 0)
				resultl, _ := models.GetAllSvginfo(query, []string{}, []string{"s_id"}, []string{"asc"}, rand, 1)

				for _, v := range resultl {
					svg += v.(models.Svg_info).Svg_dtl
					// link the next svg to be strcat
					if v.(models.Svg_info).N_id != 0 {
						query := make(map[string]string, 0)
						query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
						query["s_id"] = fmt.Sprintf("%v", v.(models.Svg_info).N_id)
						models.GetAllSvginfo(query, []string{}, []string{}, []string{}, 0, 1)
						for _, v := range resultl {
							svg += v.(models.Svg_info).Svg_dtl
						}
					}
				}
			} else {
				// bodyline
				count := models.GetCountByCatagoryId(catagory_id)

				if count == 0 {
					break
				}
				rand := generate_rand(count)

				query := make(map[string]string, 0)
				query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
				query["p_id"] = fmt.Sprintf("%v", 0)
				resultl, _ := models.GetAllSvginfo(query, []string{}, []string{"s_id"}, []string{"asc"}, rand, 1)

				for _, v := range resultl {
					svg += v.(models.Svg_info).Svg_dtl
					// link the next svg to be strcat
					if v.(models.Svg_info).N_id != 0 {
						query := make(map[string]string, 0)
						query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
						query["s_id"] = fmt.Sprintf("%v", v.(models.Svg_info).N_id)
						models.GetAllSvginfo(query, []string{}, []string{}, []string{}, 0, 1)
						for _, v := range resultl {
							svg += v.(models.Svg_info).Svg_dtl
						}
					}
				}
			}

		case 1:
			//use random element
			//Be or not be
			rand := generate_rand(int64(100 / percent))

			if (rand != 0 && percent != 1) || (flag == 1 && percent == 1) {
				//get count of this item
				count := models.GetCountByCatagoryId(catagory_id)
				if count == 0 {
					break
				}
				rand := generate_rand(count)

				query := make(map[string]string, 0)
				query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
				query["p_id"] = fmt.Sprintf("%v", 0)
				resultl, _ := models.GetAllSvginfo(query, []string{}, []string{"s_id"}, []string{"asc"}, rand, 1)

				for _, v := range resultl {
					svg += v.(models.Svg_info).Svg_dtl
					// link the next svg to be strcat
					if v.(models.Svg_info).N_id != 0 {
						query := make(map[string]string, 0)
						query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
						query["s_id"] = fmt.Sprintf("%v", v.(models.Svg_info).N_id)
						models.GetAllSvginfo(query, []string{}, []string{}, []string{}, 0, 1)
						for _, v := range resultl {
							svg += v.(models.Svg_info).Svg_dtl
						}
					}
				}
			}

		default:
			//choose one; drop another
			//How many items where select_flag = this, calculator only first time.
			if selectflag_times[select_flag] == selectflag_check[select_flag] {
				if bodycolor_flag == 1 {
					rand := generate_rand(models.GetCountByCatagoryId(catagory_id) / 7)
					if rand == -1 {
						break
					}

					query := make(map[string]string, 0)
					query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
					query["base_color"] = fmt.Sprintf("%v", color)
					query["p_id"] = fmt.Sprintf("%v", 0)
					resultl, _ := models.GetAllSvginfo(query, []string{}, []string{"s_id"}, []string{"asc"}, rand, 1)

					for _, v := range resultl {
						svg += v.(models.Svg_info).Svg_dtl
						if v.(models.Svg_info).N_id != 0 {
							// link the next svg to be strcat
							query := make(map[string]string, 0)
							query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
							query["s_id"] = fmt.Sprintf("%v", v.(models.Svg_info).N_id)
							models.GetAllSvginfo(query, []string{}, []string{}, []string{}, 0, 1)
							for _, v := range resultl {
								svg += v.(models.Svg_info).Svg_dtl
							}
						}
					}
				} else {
					rand := generate_rand(models.GetCountByCatagoryId(catagory_id))
					if rand == -1 {
						break
					}

					query := make(map[string]string, 0)
					query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
					query["p_id"] = fmt.Sprintf("%v", 0)
					resultl, _ := models.GetAllSvginfo(query, []string{}, []string{"s_id"}, []string{"asc"}, rand, 1)
					for _, v := range resultl {
						svg += v.(models.Svg_info).Svg_dtl
						// link the next svg to be strcat
						if v.(models.Svg_info).N_id != 0 {
							query := make(map[string]string, 0)
							query["catagory_id"] = fmt.Sprintf("%v", catagory_id)
							query["s_id"] = fmt.Sprintf("%v", v.(models.Svg_info).N_id)
							models.GetAllSvginfo(query, []string{}, []string{}, []string{}, 0, 1)
							for _, v := range resultl {
								svg += v.(models.Svg_info).Svg_dtl
							}
						}
					}
				}
			} else {
				//do nothing
			}
			selectflag_times[select_flag]++
		}
	}

	svg += "</svg>"

	// Create svg file and write inside
	//	if basePath[0:1] != "\\" {
	//		basePath += "\\"
	//	}
	fileName := fmt.Sprintf("%s%v.svg", time.Now().Format("20060102150405"), petID)
	strFile := basePath + fileName
	f, err := os.Create(strFile)
	w := bufio.NewWriter(f)
	if _, err = w.WriteString(svg); err != nil {
		//err deal
		fmt.Println(err)
	}
	w.Flush()

	f.Close()

	return fileName
}
