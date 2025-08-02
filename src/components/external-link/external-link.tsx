import { type AnchorHTMLAttributes } from "react";

type ExternalLinkProps = AnchorHTMLAttributes<HTMLAnchorElement>;

export const ExternalLink = ({
  ref,
  ...props
}: ExternalLinkProps & { ref?: React.RefObject<HTMLAnchorElement | null> }) => (
  <a target="_blank" rel="noopener noreferrer" {...props} ref={ref} />
);
