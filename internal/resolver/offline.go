package resolver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type OfflineResolver struct{}

func NewOffline() *OfflineResolver {
	return &OfflineResolver{}
}

type npmResult struct {
	Objects []struct {
		Package struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"package"`
	} `json:"objects"`
}

type nugetResult struct {
	Data []struct {
		ID          string `json:"id"`
		Description string `json:"description"`
	} `json:"data"`
}

type cargoResult struct {
	Crates []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"crates"`
}

func (r *OfflineResolver) Resolve(language, description string) ([]Result, error) {
	lang := strings.ToLower(language)
	switch lang {
	case "javascript", "js", "node", "typescript", "ts":
		return r.searchNPM(description)
	case "rust", "rs":
		return r.searchCargo(description)
	case "c#", "csharp", "dotnet", ".net", "nuget":
		return r.searchNuGet(description)
	case "python", "py":
		return r.searchNPM(description)
	case "go", "golang":
		return r.searchNPM(description)
	default:
		return r.searchNPM(description)
	}
}

func (r *OfflineResolver) searchNPM(desc string) ([]Result, error) {
	u := fmt.Sprintf("https://registry.npmjs.org/-/v1/search?text=%s&size=10", url.QueryEscape(desc))
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("npm search: %w", err)
	}
	defer resp.Body.Close()

	var data npmResult
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("npm decode: %w", err)
	}

	results := make([]Result, 0, len(data.Objects))
	for _, obj := range data.Objects {
		results = append(results, Result{
			Library: Library{
				Name:        obj.Package.Name,
				Description: obj.Package.Description,
				ImportPath:  obj.Package.Name,
				Language:    "javascript",
			},
		})
	}
	return results, nil
}

func (r *OfflineResolver) searchNuGet(desc string) ([]Result, error) {
	u := fmt.Sprintf("https://api-v2v3search-0.nuget.org/query?q=%s&take=10", url.QueryEscape(desc))
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("nuget search: %w", err)
	}
	defer resp.Body.Close()

	var data nugetResult
	body, _ := io.ReadAll(resp.Body)
	body = bytes.TrimPrefix(body, []byte{0xEF, 0xBB, 0xBF})
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("nuget decode: %w", err)
	}

	results := make([]Result, 0, len(data.Data))
	for _, d := range data.Data {
		results = append(results, Result{
			Library: Library{
				Name:        d.ID,
				Description: d.Description,
				ImportPath:  d.ID,
				Language:    "c#",
			},
		})
	}
	return results, nil
}

func (r *OfflineResolver) searchCargo(desc string) ([]Result, error) {
	u := fmt.Sprintf("https://crates.io/api/v1/crates?q=%s&per_page=10", url.QueryEscape(desc))
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "findlib/0.1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cargo search: %w", err)
	}
	defer resp.Body.Close()

	var data cargoResult
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("cargo decode: %w", err)
	}

	results := make([]Result, 0, len(data.Crates))
	for _, c := range data.Crates {
		results = append(results, Result{
			Library: Library{
				Name:        c.Name,
				Description: c.Description,
				ImportPath:  c.Name,
				Language:    "rust",
			},
		})
	}
	return results, nil
}
