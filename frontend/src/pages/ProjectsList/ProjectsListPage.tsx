import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useProjectStore } from '../../store/useProjectStore';
import { projectsApi } from '../../api/projects';
import { calculationApi } from '../../api/calculation';
import { formatCOP } from '../../utils/format';
import { Project } from '../../types';
import Card from '../../components/common/Card';
import Button from '../../components/common/Button';

export default function ProjectsListPage() {
  const navigate = useNavigate();
  const {
    projects,
    setProjects,
    setCurrentProject,
    setCurrentScenario,
    setScenarios,
    setIsLoading,
  } = useProjectStore();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    projectsApi
      .getAll()
      .then(setProjects)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  const handleOpenProject = async (project: Project) => {
    setCurrentProject(project);
    setIsLoading(true);
    try {
      const scenarios = await projectsApi.getScenarios(project._id);
      setScenarios(scenarios);
      if (scenarios.length > 0) {
        setCurrentScenario(scenarios[0]);
        navigate('/results');
      } else {
        navigate('/new');
      }
    } catch (err) {
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Estas seguro de eliminar este proyecto?')) return;
    try {
      await projectsApi.delete(id);
      setProjects(projects.filter((p) => p._id !== id));
    } catch (err) {
      console.error(err);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="w-8 h-8 border-4 border-brand border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-fg-primary">Mis Proyectos</h2>
        <Button onClick={() => navigate('/new')} size="lg">
          + Nuevo Proyecto
        </Button>
      </div>

      {projects.length === 0 ? (
        <Card padding="lg" className="text-center py-20">
          <p className="text-fg-muted text-lg mb-2">No tienes proyectos aun.</p>
          <p className="text-fg-muted text-sm mb-6">
            Crea tu primer proyecto para dimensionar un sistema solar.
          </p>
          <Button onClick={() => navigate('/new')}>
            Crear Proyecto
          </Button>
        </Card>
      ) : (
        <div className="space-y-3">
          {projects.map((project) => (
            <Card
              key={project._id}
              hoverable
              padding="sm"
              onClick={() => handleOpenProject(project)}
            >
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-semibold text-fg-primary">{project.name}</h3>
                  <p className="text-sm text-fg-secondary">
                    {project.location.city}, {project.location.department} |{' '}
                    <span className="font-mono tabular-nums">{project.location.latitude.toFixed(2)}</span>°,{' '}
                    <span className="font-mono tabular-nums">{project.location.longitude.toFixed(2)}</span>°
                  </p>
                  <div className="flex gap-4 mt-2 text-xs text-fg-tertiary">
                    <span>Estrato <span className="font-mono tabular-nums">{project.consumption.estrato}</span></span>
                    <span>Techo: <span className="font-mono tabular-nums">{project.roof.area}</span>m2</span>
                    <span>Tipo: {project.systemType}</span>
                    <span><span className="font-mono tabular-nums">{project.scenarios?.length || 0}</span> escenarios</span>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-xs text-fg-tertiary">
                    {new Date(project.updatedAt).toLocaleDateString('es-CO')}
                  </span>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="hover:text-danger hover:bg-danger-soft"
                    onClick={(e) => {
                      e.stopPropagation();
                      handleDelete(project._id);
                    }}
                  >
                    Eliminar
                  </Button>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
