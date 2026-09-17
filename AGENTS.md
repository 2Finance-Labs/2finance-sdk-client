# 2finance-sdk-client — instruções de trabalho

<!-- 2finance-hub:start -->
## Arquitetura e planejamento do ecossistema 2Finance

Antes de trabalho neste componente, leia a [arquitetura central](../2finance-docs/standard/docs/architecture/system-overview.mdx), a entrada deste repositório no [mapa de componentes](../2finance-docs/standard/docs/architecture/repository-map.md) e os WPs afetados no [plano master](../2finance-docs/standard/plans/PLAN-2FINANCE-0001-master.md).

O hub coordena arquitetura, contratos compartilhados e planos transversais. Detalhes de implementação permanecem neste repositório. Mudanças de componentes, dependências ou contratos devem ser desenhadas no hub antes da implementação e reconciliadas ao final; integrações seguem schema-first.

Ao alterar ou criar plano local, mantenha o vínculo de volta, registre a fonte e reconcilie o master, estado, evidências e hash revisado em `2finance-docs/standard/plans/registry.yaml`. Não concluir a tarefa apenas atualizando checkboxes ou hashes. Valide com `npm run check:workspace` no hub quando os checkouts estiverem disponíveis; indisponibilidade deve ser registrada, sem alegar validação integral.

Estas regras não substituem as instruções específicas do componente nem ampliam autorização para mudanças operacionais, externas ou destrutivas. Siga também a governança MYREPO quando estiver no workspace.
<!-- 2finance-hub:end -->
