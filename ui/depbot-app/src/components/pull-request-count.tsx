import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Typography from '@mui/material/Typography';
import CardActionArea from '@mui/material/CardActionArea';

type PullRequestCountProps = {
  title: string;
  count: number;
}

const PullRequestCount: React.FC<PullRequestCountProps> = ({ title, count }: PullRequestCountProps) => {
  return (
    <Card sx={{ paddingLeft: '2.5vh', paddingRight: '2.5vh' }}>
      <CardActionArea>
        <CardContent>
          <Typography gutterBottom variant="h5" component="div">
            {count}
          </Typography>
          <Typography variant="body2" sx={{ color: 'text.secondary' }}>
            {title}
          </Typography>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}

export default PullRequestCount;

