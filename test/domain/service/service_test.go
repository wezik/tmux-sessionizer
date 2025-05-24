package service_test

import (
	"fmt"
	. "phopper/src/domain/model"
	. "phopper/src/domain/service"
	"testing"
)

func Test_Service(t *testing.T) {
	t.Run("create project", func(t *testing.T) {
		t.Run("creates project", func(t *testing.T) {
			// given
			st := &MockStorage{}
			svc := NewService(nil, nil, st)
			cwd := "/home/test"
			name := "foobar"

			// when
			svc.CreateProject(cwd, name)

			// then
			if st.SaveCalls != 1 {
				t.Errorf("Save should be called once")
			}

			if st.SaveParam1.Name != name {
				t.Errorf("Saved project should have name %s", name)
			}

			if st.SaveParam1.Template.Root != cwd {
				t.Errorf("Saved project should have root %s", cwd)
			}

			if st.SaveParam1.Template.Name != name {
				t.Errorf("Saved project should have template name %s", name)
			}

			if len(st.SaveParam1.Template.Windows) != 1 {
				t.Errorf("Saved project should have one window")
			}

			if st.SaveParam1.Template.Windows[0].Name != "shell" {
				t.Errorf("Saved project should have window name %s", "shell")
			}
		})

		t.Run("panics with invalid data", func(t *testing.T) {
			for _, args := range [][]string{
				{"", ""},
				{"", "foo"},
				{"/foo/bar", ""},
			} {
				t.Run(fmt.Sprintf("for %s and %s", args[0], args[1]), func(t *testing.T) {
					// given
					svc := NewService(nil, nil, nil)
					cwd := args[0]
					name := args[1]

					// expect
					defer func() {
						if r := recover(); r == nil {
							t.Errorf("The code did not panic")
						}
					}()
					svc.CreateProject(cwd, name)
				})
			}
		})
	})

	t.Run("select and open project", func(t *testing.T) {
		t.Run("finds and opens project with multiplexer", func(t *testing.T) {
			// given
			name := "foobar"
			sc := &MockSelector{}
			mu := &MockMultiplexer{}
			p := &Project{ID: "1234", Name: name}
			st := &MockStorage{}
			st.FindReturn = p
			svc := NewService(sc, mu, st)

			// when
			svc.SelectAndOpenProject(name)

			// then
			if st.FindCalls != 1 {
				t.Errorf("Find should be called once")
			}

			if st.FindParam1 != name {
				t.Errorf("Find should be called with %s", name)
			}

			if sc.SelectFromCalls != 0 {
				// shouldn't start selector
				t.Errorf("Selector should not be called")
			}

			if mu.AttachProjectCalls != 1 {
				t.Errorf("The project should be attached")
			}

			if mu.AttachProjectParam1.ID != p.ID {
				t.Errorf("The project should be attached with %s", p.ID)
			}
		})

		t.Run("selects from selector and opens project with multiplexer", func(t *testing.T) {
			// given
			name := "foobar"
			p := &Project{ID: "1234", Name: name}
			pjs := []*Project{p}
			pjsStringified := []string{p.Name}

			sc := &MockSelector{}
			sc.SelectFromReturn = p.Name

			mu := &MockMultiplexer{}

			st := &MockStorage{}
			st.ListReturn = pjs

			svc := NewService(sc, mu, st)

			// when
			svc.SelectAndOpenProject("")

			// then
			if st.FindCalls != 0 {
				t.Errorf("Find should not be called")
			}

			if sc.SelectFromCalls != 1 {
				t.Errorf("Selector should be called once")
			}

			for i, pj := range pjsStringified {
				if i >= len(sc.SelectFromParam1) {
					t.Errorf("Selector should be called with %s got %s", pjsStringified, sc.SelectFromParam1)
				}

				if pj != pjsStringified[i] {
					t.Errorf("Selector should be called with %s got %s", pjsStringified, sc.SelectFromParam1)
				}
			}

			if mu.AttachProjectCalls != 1 {
				t.Errorf("The project should be attached")
			}

			if mu.AttachProjectParam1 != p {
				t.Errorf("The project should be attached with %s got %s", p.Name, mu.AttachProjectParam1.Name)
			}
		})

		t.Run("panics when project is not found", func(t *testing.T) {
			// given
			name := "foobar"
			sc := &MockSelector{}
			mu := &MockMultiplexer{}
			err := ErrNotFound
			st := &MockStorage{FindErr: err}
			svc := NewService(sc, mu, st)

			// when
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("The code did not panic")
				} else if r.(error) != ErrNotFound {
					t.Errorf("The error should be %s was %s", err, r)
				}
			}()
			svc.SelectAndOpenProject(name)

			// then
			if st.FindCalls != 1 {
				t.Errorf("Find should be called once")
			}

			if st.FindParam1 != name {
				t.Errorf("Find should be called with %s", name)
			}

			if sc.SelectFromCalls != 0 {
				t.Errorf("Selector should not be called")
			}

			if mu.AttachProjectCalls != 0 {
				t.Errorf("The project should not be attached")
			}
		})
	})
}
